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
  },
  plugins: [],
}
