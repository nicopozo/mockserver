module.exports = {
    configureWebpack: {
        devtool: 'source-map'
    },
    publicPath: process.env.NODE_ENV === 'production'
        ? '/mock-service/admin/'
        : '/admin/',

    transpileDependencies: [
            'vuetify'
          ]
}