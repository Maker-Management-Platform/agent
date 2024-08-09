/** @type {import('tailwindcss').Config} */
export default {
  content: [
    '../**/*.templ',
    './src/**/*.js'],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
}