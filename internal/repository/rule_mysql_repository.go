package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	guuid "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
)

type RuleMySQLRepository struct {
	db Database
}

func NewRuleMySQLRepository(db Database) IRuleRepository {
	return &RuleMySQLRepository{
		db: db,
	}
}

type RuleRow struct {
	Key         string `db:"key"`
	Application string `db:"application"`
	Name        string `db:"name"`
	Path        string `db:"path"`
	Strategy    string `db:"strategy"`
	Method      string `db:"method"`
	Status      string `db:"status"`
	Pattern     string `db:"pattern"`
}

type VariableRow struct {
	ID      int64  `db:"id"`
	Type    string `db:"type"`
	Name    string `db:"name"`
	Key     string `db:"key"`
	RuleKey string `db:"rule_key"`
}

type ResponseRow struct {
	ID          int64   `db:"id"`
	Body        string  `db:"body"`
	ContentType string  `db:"content_type"`
	HTTPStatus  int     `db:"http_status"`
	Delay       int     `db:"delay"`
	Scene       *string `db:"scene"`
	RuleKey     string  `db:"rule_key"`
}

func (repository *RuleMySQLRepository) Create(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	query := "INSERT INTO rules (`key`, application, name, path, strategy, method, status, pattern) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := repository.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error strating transaction, %w", err)
	}

	defer repository.commitOrRollback(ctx, tx, err)

	rule.Key = fmt.Sprintf("%v", guuid.New())

	_, err = tx.Exec(query, rule.Key, rule.Application, rule.Name, rule.Path, rule.Strategy, rule.Method, rule.Status,
		CreateExpression(rule.Path))

	if err != nil {
		logger.Error(repository, nil, err, "error creating rule in DB")

		return nil, fmt.Errorf("error creating rule into DB, %w", err)
	}

	err = repository.insertVariables(ctx, rule, tx)
	if err != nil {
		return nil, err
	}

	err = repository.insertResponses(ctx, rule, tx)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (repository *RuleMySQLRepository) Update(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	query := "UPDATE rules SET application=?, name=?, path=?, strategy=?, method=?, status=?, pattern=? WHERE `key`=?"

	tx, err := repository.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction, %w", err)
	}

	defer repository.commitOrRollback(ctx, tx, err)

	_, err = tx.Exec(query, rule.Application, rule.Name, rule.Path, rule.Strategy, rule.Method, rule.Status,
		CreateExpression(rule.Path), rule.Key)
	if err != nil {
		logger.Error(repository, nil, err, "error updating rule in DB")

		return nil, fmt.Errorf("error updating item into DB, %w", err)
	}

	err = repository.deleteVariables(ctx, rule.Key, tx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	err = repository.insertVariables(ctx, rule, tx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	err = repository.deleteResponses(ctx, rule.Key, tx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	err = repository.insertResponses(ctx, rule, tx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	return rule, nil
}

func (repository *RuleMySQLRepository) deleteResponses(ctx context.Context, key string, tx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := "DELETE FROM responses WHERE rule_key=?"

	_, err := tx.Exec(query, key)

	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return fmt.Errorf("error updating task, %w", err)
	}

	return nil
}

func (repository *RuleMySQLRepository) deleteVariables(ctx context.Context, key string, tx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := "DELETE FROM variables WHERE rule_key=?"

	_, err := tx.Exec(query, key)

	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return fmt.Errorf("error updating task, %w", err)
	}

	return nil
}

func (repository *RuleMySQLRepository) Get(ctx context.Context, key string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	query := "SELECT * FROM rules WHERE `key` = ?"
	row := RuleRow{}

	err = repository.db.Get(&row, query, key)

	if err != nil {
		if err.Error() == noRowsMessage {
			msg := fmt.Sprintf("no rule found with key: %s", key)
			logger.Error(repository, nil, err, msg)

			return nil, mockserrors.RuleNotFoundError{
				Message: msg,
			}
		}

		logger.Error(repository, nil, err, "error getting rule from DB")

		return nil, fmt.Errorf("error getting rule from DB, %w", err)
	}

	var variables []VariableRow

	query = `SELECT * FROM variables
				WHERE rule_key = ?`

	err = repository.db.Select(&variables, query, key)

	if err != nil {
		logger.Error(repository, nil, err, "error executing SQL query")

		return nil, fmt.Errorf("error getting variables for rule with key '%s' from DB, %w", key, err)
	}

	var responses []ResponseRow

	query = `SELECT * FROM responses
				WHERE rule_key = ?`

	err = repository.db.Select(&responses, query, key)

	if err != nil {
		logger.Error(repository, nil, err, "error executing SQL query")

		return nil, fmt.Errorf("error getting responses for rule with key '%s' from DB, %w", key, err)
	}

	return parseRule(row, variables, responses), nil
}

func (repository *RuleMySQLRepository) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	var rows []RuleRow

	var err error

	searchQuery := newSearchQuery(params)

	err = repository.db.Select(&rows, searchQuery, paging.Limit, paging.Offset)

	if err != nil {
		logger.Error(repository, nil, err, "error executing SQL query")

		return nil, fmt.Errorf("error searching rules in DB, %w", err)
	}

	if len(rows) > 0 {
		var total int64

		totalQuery := "SELECT COUNT(*) as total FROM rules " + newWhereClause(params)

		err = repository.db.Get(&total, totalQuery)

		if err != nil {
			logger.Error(repository, nil, err, "error executing SQL query")

			return nil, fmt.Errorf("error calculating total rules in DB, %w", err)
		}

		paging.Total = total
	}

	rules := make([]*model.Rule, len(rows))
	for index, row := range rows {
		rules[index], err = repository.Get(ctx, row.Key)

		if err != nil {
			return nil, err
		}
	}

	return &model.RuleList{Paging: paging, Results: rules}, nil
}

func (repository *RuleMySQLRepository) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	var err error

	tx, err := repository.db.Beginx()
	if err != nil {
		return fmt.Errorf("error strating transaction, %w", err)
	}

	defer repository.commitOrRollback(ctx, tx, err)

	err = repository.deleteVariables(ctx, key, tx)
	if err != nil {
		logger.Error(repository, nil, err, "error deleting task in DB")

		return fmt.Errorf("error deleting task, %w", err)
	}

	err = repository.deleteResponses(ctx, key, tx)
	if err != nil {
		logger.Error(repository, nil, err, "error deleting task in DB")

		return fmt.Errorf("error deleting task, %w", err)
	}

	query := "DELETE FROM rules WHERE `key`=?"

	_, err = tx.Exec(query, key)

	if err != nil {
		logger.Error(repository, nil, err, "error deleting rule in DB")

		return fmt.Errorf("error deleting rule, %w", err)
	}

	return nil
}

func (repository *RuleMySQLRepository) SearchByMethodAndPath(ctx context.Context, method string,
	path string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	var rows []RuleRow

	searchQuery := "SELECT `key`, pattern, status FROM rules WHERE method = ?"

	err = repository.db.Select(&rows, searchQuery, strings.ToUpper(method))

	if err != nil {
		logger.Error(repository, nil, err, "error executing SQL query")

		return nil, fmt.Errorf("error searching rules in DB, %w", err)
	}

	for _, row := range rows {
		var regex = regexp.MustCompile(row.Pattern)

		if regex.MatchString(path) {
			if row.Status == model.RuleStatusEnabled {
				return repository.Get(ctx, row.Key)
			}
		}
	}

	return nil, mockserrors.RuleNotFoundError{
		Message: fmt.Sprintf("no rule found for path: %s and method %s", path, method),
	}
}

func (repository *RuleMySQLRepository) insertVariables(ctx context.Context, rule *model.Rule, tx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := "INSERT INTO variables (type, name, `key`, rule_key) VALUES (?, ?, ?, ?)"

	for _, v := range rule.Variables {
		_, err := tx.Exec(query, v.Type, v.Name, v.Key, rule.Key)

		if err != nil {
			logger.Error(repository, nil, err, "error creating rule variable in DB")

			return fmt.Errorf("error creating rule variable in DB, %w", err)
		}
	}

	return nil
}

func (repository *RuleMySQLRepository) insertResponses(ctx context.Context, rule *model.Rule, tx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := `INSERT INTO responses (body, content_type, http_status, delay, scene, rule_key) VALUES (?, ?, ?, ?, ?, ?)`

	for _, r := range rule.Responses {
		_, err := tx.Exec(query, r.Body, r.ContentType, r.HTTPStatus, r.Delay, r.Scene, rule.Key)

		if err != nil {
			logger.Error(repository, nil, err, "error creating rule response in DB")

			return fmt.Errorf("error creating rule response in DB, %w", err)
		}
	}

	return nil
}

func newSearchQuery(params map[string]interface{}) string {
	query := "SELECT * FROM rules"
	where := newWhereClause(params)
	order := " ORDER BY application, path, method LIMIT ? OFFSET ?"

	return query + where + order
}

func newWhereClause(params map[string]interface{}) string {
	if len(params) == 0 {
		return " "
	}

	where := " WHERE "

	index := 0
	for key, value := range params {
		if index > 0 {
			where += " AND "
		}

		switch key {
		case "application", "status", "method", "pattern", "strategy", "path", "name", "key":
			v := strings.ToLower(fmt.Sprintf("%v", value))
			where += key + " like '%" + v + "%'"
		default:
			where += key + "=" + fmt.Sprintf("%v", value)
		}
		index++
	}

	return where
}

func parseRule(row RuleRow, variables []VariableRow, responses []ResponseRow) *model.Rule {
	vars := make([]*model.Variable, 0)

	for _, v := range variables {
		newVar := model.Variable{
			Type: v.Type,
			Name: v.Name,
			Key:  v.Key,
		}

		vars = append(vars, &newVar)
	}

	resps := make([]model.Response, 0)

	for _, r := range responses {
		scene := ""

		if r.Scene != nil {
			scene = *r.Scene
		}

		newResp := model.Response{
			Body:        r.Body,
			ContentType: r.ContentType,
			HTTPStatus:  r.HTTPStatus,
			Delay:       r.Delay,
			Scene:       scene,
		}

		resps = append(resps, newResp)
	}

	return &model.Rule{
		Key:         row.Key,
		Application: row.Application,
		Name:        row.Name,
		Path:        row.Path,
		Strategy:    row.Strategy,
		Method:      row.Method,
		Status:      row.Status,
		Variables:   vars,
		Responses:   resps,
	}
}

func (repository *RuleMySQLRepository) commitOrRollback(ctx context.Context, tx *sqlx.Tx, err error) {
	logger := mockscontext.Logger(ctx)

	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			logger.Error(repository, nil, err, "Error committing changes")
			_ = tx.Rollback()
		}
	}
}
