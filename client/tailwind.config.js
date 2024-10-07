/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'primary': '#191825',
        'secondary': '#865DFF',
        'accent-1': '#E384FF',
        'accent-2': '#FFA3FD',
        'rgb-secondary': 'rgb(134, 93, 255)'
      },
      keyframes: {
        'border-spin': {
          '100%': {
            transform: 'rotate(360deg)',
          }
        }
      },
      animation: {
        'border-spin': 'border-spin 3s linear infinite',
      }
    },
  },
  plugins: [],
}
