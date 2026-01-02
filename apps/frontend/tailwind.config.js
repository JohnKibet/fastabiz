module.exports = {
  content: ["./Pages/**/*.{razor,html,cshtml}", "./Pages/Shared/**/*.{razor,html,cshtml}"],
  safelist: ['no-underline'],
  theme: {
    extend: {
      colors: {
      primary: "#2563EB",   // Blue
      accent: "#10B981",    // Emerald
      ctaHover: "#059669",  // CTA hover
      bgSoft: "#F9FAFB",    // Background
      textDark: "#111827",  // Headings
      textLight: "#6B7280", // Subtext
      stoneblue: {
        DEFAULT: '#59788E',  // your main shade
        light: '#E5EEF3',    // soft background tint
        hover: '#F5F8FA'     // subtle hover background
      },
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui'],
        poppins: ['Poppins', 'sans-serif'],
      },
    },
    animation: {
    'slide-down': 'slideDown 0.2s ease-out',
    'drawer-in': 'drawerIn .35s ease-out forwards',
    'drawer-out': 'drawerOut .3s ease-in forwards',
    'drawer-secondary-in': 'drawerSecondaryIn .35s ease-out forwards',
    'drawer-secondary-out': 'drawerSecondaryOut .3s ease-in forwards',
    },
    keyframes: {
      slideDown: {
        '0%': { opacity: 0, transform: 'translateY(-10px)' },
        '100%': { opacity: 1, transform: 'translateY(0)' }
      },
      drawerIn: {
      '0%': { transform: 'translateX(-100%)' },
      '100%': { transform: 'translateX(0)' }
      },
      drawerOut: {
        '0%': { transform: 'translateX(0)' },
        '100%': { transform: 'translateX(100%)' }
      },
      drawerSecondaryIn: {
        '0%': { transform: 'translateX(-100%)' },
        '100%': { transform: 'translateX(0)' }
      },
      drawerSecondaryOut: {
        '0%': { transform: 'translateX(0)' },
        '100%': { transform: 'translateX(-100%)' }
      },
      fadeIn: {
        '0%': { opacity: 0 },
        '100%': { opacity: 1 }
      },
      fadeOut: {
        '0%': { opacity: 1 },
        '100%': { opacity: 0 }
      }
    }
  },
  plugins: [],
}
