import { AuthState } from 'models/auth.models';
import { createContext } from 'react';

export const AuthContext = createContext<AuthState>({} as AuthState);
