import { AuthContext } from 'context/auth-provider/AuthContext';
import { AuthState } from 'models/auth.models';
import { useContext } from 'react';

export default function useAuthState(): AuthState {
	return useContext(AuthContext);
}
