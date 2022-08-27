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
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

type ruleMySQLRepository struct {
	db Database
}

func NewRuleMySQLRepository(db Database) RuleRepository {
	return &ruleMySQLRepository{
		db: db,
	}
}

type RuleRow struct {
	Key               string `db:"key"`
	Group             string `db:"group"`
	Name              string `db:"name"`
	Path              string `db:"path"`
	Strategy          string `db:"strategy"`
	Method            string `db:"method"`
	Status            string `db:"status"`
	Pattern           string `db:"pattern"`
	NextResponseIndex int    `db:"next_response_index"`
}

type VariableRow struct {
	ID         int64   `db:"id"`
	Type       string  `db:"type"`
	Name       string  `db:"name"`
	Key        string  `db:"key"`
	RuleKey    string  `db:"rule_key"`
	Assertions *string `db:"assertions"`
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

func (repository *ruleMySQLRepository) Create(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	query := "INSERT INTO rules (`key`, `group`, name, path, strategy, method, status, pattern, next_response_index) " +
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	trx, err := repository.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error strating transaction, %w", err)
	}

	defer repository.commitOrRollback(ctx, trx, err)

	rule.Key = fmt.Sprintf("%v", guuid.New())

	_, err = trx.Exec(query, rule.Key, rule.Group, rule.Name, rule.Path, rule.Strategy, rule.Method, rule.Status,
		CreateExpression(rule.Path), rule.NextResponseIndex)

	if err != nil {
		logger.Error(repository, nil, err, "error creating rule in DB")

		return nil, fmt.Errorf("error creating rule into DB, %w", err)
	}

	err = repository.insertVariables(ctx, rule, trx)
	if err != nil {
		return nil, err
	}

	err = repository.insertResponses(ctx, rule, trx)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (repository *ruleMySQLRepository) Update(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	query := "UPDATE rules SET `group`=?, name=?, path=?, strategy=?, method=?, status=?, pattern=?, " +
		" next_response_index=? WHERE `key`=?"

	trx, err := repository.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction, %w", err)
	}

	defer repository.commitOrRollback(ctx, trx, err)

	_, err = trx.Exec(query, rule.Group, rule.Name, rule.Path, rule.Strategy, rule.Method, rule.Status,
		CreateExpression(rule.Path), rule.NextResponseIndex, rule.Key)
	if err != nil {
		logger.Error(repository, nil, err, "error updating rule in DB")

		return nil, fmt.Errorf("error updating item into DB, %w", err)
	}

	err = repository.deleteVariables(ctx, rule.Key, trx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	err = repository.insertVariables(ctx, rule, trx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	err = repository.deleteResponses(ctx, rule.Key, trx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	err = repository.insertResponses(ctx, rule, trx)
	if err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return nil, fmt.Errorf("error updating task, %w", err)
	}

	return rule, nil
}

func (repository *ruleMySQLRepository) deleteResponses(ctx context.Context, key string, tx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := "DELETE FROM responses WHERE rule_key=?"

	if _, err := tx.Exec(query, key); err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return fmt.Errorf("error updating task, %w", err)
	}

	return nil
}

func (repository *ruleMySQLRepository) deleteVariables(ctx context.Context, key string, tx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := "DELETE FROM variables WHERE rule_key=?"

	if _, err := tx.Exec(query, key); err != nil {
		logger.Error(repository, nil, err, "error updating task in DB")

		return fmt.Errorf("error updating task, %w", err)
	}

	return nil
}

func (repository *ruleMySQLRepository) Get(ctx context.Context, key string) (*model.Rule, error) {
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

	query = "SELECT id, type, name, `key`, rule_key, assertions FROM variables WHERE rule_key = ?"

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

func (repository *ruleMySQLRepository) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging,
) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	var rows []RuleRow

	var err error

	searchQuery, err := newSearchQuery(params)
	if err != nil {
		return nil, err
	}

	err = repository.db.Select(&rows, searchQuery, paging.Limit, paging.Offset)

	if err != nil {
		logger.Error(repository, nil, err, "error executing SQL query")

		return nil, fmt.Errorf("error searching rules in DB, %w", err)
	}

	if len(rows) > 0 {
		var total int64

		where, err := newWhereClause(params)
		if err != nil {
			return nil, err
		}

		totalQuery := "SELECT COUNT(*) as total FROM rules " + where

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

func (repository *ruleMySQLRepository) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	var err error

	trx, err := repository.db.Beginx()
	if err != nil {
		return fmt.Errorf("error strating transaction, %w", err)
	}

	defer repository.commitOrRollback(ctx, trx, err)

	err = repository.deleteVariables(ctx, key, trx)
	if err != nil {
		logger.Error(repository, nil, err, "error deleting task in DB")

		return fmt.Errorf("error deleting task, %w", err)
	}

	err = repository.deleteResponses(ctx, key, trx)
	if err != nil {
		logger.Error(repository, nil, err, "error deleting task in DB")

		return fmt.Errorf("error deleting task, %w", err)
	}

	query := "DELETE FROM rules WHERE `key`=?"

	_, err = trx.Exec(query, key)

	if err != nil {
		logger.Error(repository, nil, err, "error deleting rule in DB")

		return fmt.Errorf("error deleting rule, %w", err)
	}

	return nil
}

func (repository *ruleMySQLRepository) SearchByMethodAndPath(ctx context.Context, method string,
	path string,
) (*model.Rule, error) {
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
		regex := regexp.MustCompile(row.Pattern)

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

func (repository *ruleMySQLRepository) insertVariables(ctx context.Context, rule *model.Rule, trx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := "INSERT INTO variables (type, name, `key`, rule_key, assertions) VALUES (?, ?, ?, ?, ?)"

	for _, variable := range rule.Variables {
		var assertions *string

		if variable.Assertions != nil {
			a := jsonutils.Marshal(variable.Assertions)
			assertions = &a
		}

		_, err := trx.Exec(query, variable.Type, variable.Name, variable.Key, rule.Key, assertions)
		if err != nil {
			logger.Error(repository, nil, err, "error creating rule variable in DB")

			return fmt.Errorf("error creating rule variable in DB, %w", err)
		}
	}

	return nil
}

func (repository *ruleMySQLRepository) insertResponses(ctx context.Context, rule *model.Rule, trx *sqlx.Tx) error {
	logger := mockscontext.Logger(ctx)

	query := `INSERT INTO responses (body, content_type, http_status, delay, scene, rule_key) VALUES (?, ?, ?, ?, ?, ?)`

	for _, r := range rule.Responses {
		_, err := trx.Exec(query, r.Body, r.ContentType, r.HTTPStatus, r.Delay, r.Scene, rule.Key)
		if err != nil {
			logger.Error(repository, nil, err, "error creating rule response in DB")

			return fmt.Errorf("error creating rule response in DB, %w", err)
		}
	}

	return nil
}

func newSearchQuery(params map[string]interface{}) (string, error) {
	query := "SELECT * FROM rules"

	where, err := newWhereClause(params)
	if err != nil {
		return "", err
	}

	order := " ORDER BY `group`, path, method LIMIT ? OFFSET ?"

	return query + where + order, nil
}

func newWhereClause(params map[string]interface{}) (string, error) {
	if len(params) == 0 {
		return " ", nil
	}

	where := " WHERE "

	index := 0
	for key, value := range params {
		if index > 0 {
			where += " AND "
		}

		switch key {
		case "status", "method", "pattern", "strategy", "path", "name":
			v := strings.ToLower(fmt.Sprintf("%v", value))
			where += key + " like '%" + v + "%'"
		case "group", "key":
			v := strings.ToLower(fmt.Sprintf("%v", value))
			where += "`" + key + "` like '%" + v + "%'"
		default:
			return "", mockserrors.InvalidRulesError{Message: fmt.Sprintf("%s is not a valid parameter", key)}
		}
		index++
	}

	return where, nil
}

func parseRule(row RuleRow, variables []VariableRow, responses []ResponseRow) *model.Rule {
	vars := make([]*model.Variable, 0)

	for _, variable := range variables {
		newVar := model.Variable{
			Type: variable.Type,
			Name: variable.Name,
			Key:  variable.Key,
		}

		var assertions []*model.Assertion

		if variable.Assertions != nil {
			_ = jsonutils.Unmarshal(strings.NewReader(*variable.Assertions), &assertions)
			newVar.Assertions = assertions
		}

		vars = append(vars, &newVar)
	}

	resps := make([]model.Response, 0)

	for _, resp := range responses {
		scene := ""

		if resp.Scene != nil {
			scene = *resp.Scene
		}

		newResp := model.Response{
			Body:        resp.Body,
			ContentType: resp.ContentType,
			HTTPStatus:  resp.HTTPStatus,
			Delay:       resp.Delay,
			Scene:       scene,
		}

		resps = append(resps, newResp)
	}

	return &model.Rule{
		Key:               row.Key,
		Group:             row.Group,
		Name:              row.Name,
		Path:              row.Path,
		Strategy:          row.Strategy,
		Method:            row.Method,
		Status:            row.Status,
		Variables:         vars,
		Responses:         resps,
		NextResponseIndex: row.NextResponseIndex,
	}
}

func (repository *ruleMySQLRepository) commitOrRollback(ctx context.Context, trx *sqlx.Tx, err error) {
	logger := mockscontext.Logger(ctx)

	if err != nil {
		_ = trx.Rollback()
	} else {
		err = trx.Commit()
		if err != nil {
			logger.Error(repository, nil, err, "Error committing changes")
			_ = trx.Rollback()
		}
	}
}
