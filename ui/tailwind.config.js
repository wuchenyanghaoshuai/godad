/** @type {import('tailwindcss').Config} */

export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,vue}"],
  theme: {
    container: {
      center: true,
    },
    extend: {
      colors: {
        // GoDad 主题色彩系统
        primary: {
          50: '#fef7f0',
          100: '#fdeee0',
          200: '#fbd9bf',
          300: '#f8c194',
          400: '#f4a261',
          500: '#e76f51', // 主色调：温暖橙红色
          600: '#d85a3c',
          700: '#b84a32',
          800: '#943d30',
          900: '#78342c',
        },
        secondary: {
          50: '#f0fdfa',
          100: '#ccfbf1',
          200: '#99f6e4',
          300: '#5eead4',
          400: '#2dd4bf',
          500: '#14b8a6', // 辅助色：温和青绿色
          600: '#0d9488',
          700: '#0f766e',
          800: '#115e59',
          900: '#134e4a',
        },
        warm: {
          50: '#fffbf5',
          100: '#fff7ed',
          200: '#feedc7',
          300: '#fdedd3',
          400: '#fcd9b6',
          500: '#fdba74',
          600: '#f59e0b',
          700: '#d97706',
          800: '#b45309',
          900: '#92400e',
        },
        soft: {
          50: '#fafafa',
          100: '#f4f4f5',
          200: '#e4e4e7',
          300: '#d4d4d8',
          400: '#a1a1aa',
          500: '#71717a',
          600: '#52525b',
          700: '#3f3f46',
          800: '#27272a',
          900: '#18181b',
        }
      },
      fontFamily: {
        'chinese': ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', 'Helvetica', 'Arial', 'sans-serif'],
      },
      backgroundImage: {
        'gradient-warm': 'linear-gradient(135deg, #fef7f0 0%, #f0fdfa 100%)',
        'gradient-soft': 'linear-gradient(135deg, #fffbf5 0%, #fafafa 100%)',
        'card-gradient': 'linear-gradient(145deg, rgba(255,255,255,0.9) 0%, rgba(255,255,255,0.7) 100%)',
      },
      boxShadow: {
        'soft': '0 2px 15px -3px rgba(0, 0, 0, 0.07), 0 10px 20px -2px rgba(0, 0, 0, 0.04)',
        'gentle': '0 1px 3px rgba(0, 0, 0, 0.05), 0 1px 2px rgba(0, 0, 0, 0.1)',
        'warm': '0 4px 6px -1px rgba(231, 111, 81, 0.1), 0 2px 4px -1px rgba(231, 111, 81, 0.06)',
      },
      borderRadius: {
        'xl': '1rem',
        '2xl': '1.5rem',
        '3xl': '2rem',
      }
    },
  },
  plugins: [],
};
