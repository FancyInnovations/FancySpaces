/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import {md3} from 'vuetify/blueprints'

// Composables
import {createVuetify, type ThemeDefinition} from 'vuetify'

const darkTheme: ThemeDefinition = {
  dark: true,
  colors: {
    primary: '#FFB783',
    'on-primary': '#4F2500',
    'primary-container': '#6D390B',
    'on-primary-container': '#FFDCC5',

    secondary: '#E4BFA7',
    'on-secondary': '#422B1B',
    'secondary-container': '#5B412F',
    'on-secondary-container': '#FFDCC5',

    tertiary: '#EEBF6D',
    'on-tertiary': '#422D00',
    'tertiary-container': '#5E4200',
    'on-tertiary-container': '#FFDEA8',

    error: '#FFB4AB',
    'on-error': '#690005',
    'error-container': '#93000A',
    'on-error-container': '#FFDAD6',

    background: '#19120D',
    'on-background': '#F0DFD6',

    surface: '#19120D',
    'on-surface': '#F0DFD6',
    'surface-variant': '#52443B',
    'on-surface-variant': '#D6C3B7',

    outline: '#9F8D83',
    shadow: '#000000',

    'inverse-surface': '#F0DFD6',
    'inverse-on-surface': '#382F29',
    'inverse-primary': '#8A5021',
  },
}

export default createVuetify({
  blueprint: md3,
  theme: {
    defaultTheme: 'darkTheme',
    themes: {
      darkTheme
    }
  },
})
