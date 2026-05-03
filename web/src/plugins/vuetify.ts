// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Vuetify
import { createVuetify } from 'vuetify'

const savedTheme = localStorage.getItem('mockserver-theme')
const initialTheme = savedTheme ? savedTheme : (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')

export default createVuetify({
  theme: {
    defaultTheme: initialTheme
  }
})
