import axios from 'axios';
import { create } from 'zustand';
import { IUseTokensStore } from '../../types/interfaces';
import { baseUrl } from '../auth/useAuthStore';
import { handleError } from '../recipes/useRecipesStore';

export const useTokenStore = create<IUseTokensStore>(set => ({
	token: '',
	setToken: token => set({ token }),

	getToken: async () => {
		try {
			const response = await axios.get(`${baseUrl}/auth/tgtoken`, {
				withCredentials: true,
			});

			console.log('token res: ', response);

			if (response.status === 200) {
				set({ token: response.data.token });
			}
		} catch (err) {
			console.log('Error getting token: ', err);
			alert(`Ошибка получения токена: ${handleError(err)}`);
		}
	},
}));
