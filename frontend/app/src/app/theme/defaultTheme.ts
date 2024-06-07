import { createTheme } from '@mui/material/styles';
import colors from '@app/styles/colors.module.scss';
import breakpoints from '@app/styles/breakpoints.module.scss';
import spacing from '@app/styles/spacing.module.scss';

export const defaultTheme = createTheme({
  palette: {
    text: {
      primary: colors.black,
    },
    primary: {
      light: colors.primaryLight,
      main: colors.primary,
      dark: colors.primaryDark,
    },
    error: {
      light: colors.errorLight,
      main: colors.error,
      dark: colors.errorDark,
    },
    warning: {
      light: colors.warningLight,
      main: colors.warning,
      dark: colors.warningDark,
    },
    info: {
      light: colors.infoLight,
      main: colors.info,
      dark: colors.infoDark,
    },
    success: {
      light: colors.successLight,
      main: colors.success,
      dark: colors.successDark,
    },
    background: {
      default: colors.white,
    },
  },
  breakpoints: {
    values: {
      xs: +breakpoints.xs,
      sm: +breakpoints.sm,
      md: +breakpoints.md,
      lg: +breakpoints.lg,
      xl: +breakpoints.xl,
    },
  },
  spacing: Object.values(spacing).map((spacingItem) =>
    spacingItem.startsWith('.') ? `0${spacingItem}rem` : `${spacingItem}rem`
  ),
});
