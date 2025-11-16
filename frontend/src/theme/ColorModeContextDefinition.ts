import { createContext } from 'react';
import { PaletteMode } from '@mui/material';

export type ColorMode = PaletteMode;

export type ColorModeContextValue = {
  mode: ColorMode;
  toggleColorMode: () => void;
};

export const ColorModeContext = createContext<ColorModeContextValue | undefined>(undefined);
