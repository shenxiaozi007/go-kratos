/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Plus Jakarta Sans"', 'system-ui', 'sans-serif'],
      },
      colors: {
        primary: '#ef3957',
        'bg-light': '#f8f6f6',
        'bg-dark': '#221013',
      },
      borderRadius: {
        DEFAULT: '1rem',
        lg: '2rem',
        xl: '3rem',
      },
      backdropBlur: {
        glass: '10px',
      },
      maxWidth: {
        app: '448px',
      },
    },
  },
  plugins: [],
}
