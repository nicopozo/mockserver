<template>
  <div>
    <b-navbar toggleable="lg" type="dark" variant="info">
      <img src="../assets/mock.png" width="20" height="20" />
      <b-navbar-brand></b-navbar-brand>
      <b-navbar-brand>Mocks Server</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item
            active-class="active"
            class="nav-link"
            key="listMocks"
            :to="{name: 'ListMocks'}"
            exact="true"
          >Home</b-nav-item>
        </b-navbar-nav>

        <b-navbar-nav>
          <b-nav-item
            active-class="active"
            class="nav-link"
            key="newMock"
            :to="{name: 'NewMock'}"
            exact="true"
          >New Mock</b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
export default {
  name: "NavBar"
};
</script>