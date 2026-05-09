// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Vuetify
import { createVuetify } from 'vuetify'

const savedTheme = localStorage.getItem('mockserver-theme')
const initialTheme = savedTheme ? savedTheme : (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')

export default createVuetify({
  theme: {
    defaultTheme: initialTheme,
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#0ea5e9', // Sky 500
          secondary: '#64748b', // Slate 500
          accent: '#f43f5e', // Rose 500
          background: '#f8fafc', // Slate 50
          surface: '#ffffff',
          'on-surface': '#0f172a',
        }
      },
      dark: {
        dark: true,
        colors: {
          primary: '#38bdf8', // Sky 400
          secondary: '#94a3b8', // Slate 400
          accent: '#fb7185', // Rose 400
          background: '#020617', // Slate 950
          surface: '#0f172a', // Slate 900
          'on-surface': '#f1f5f9',
        }
      }
    }
  }
})
