import axios from 'axios';
import Cookies from 'js-cookie';
import { create } from 'zustand';
import { IUseAuthStore } from '../../types/interfaces';
import { handleError } from '../recipes/useRecipesStore';

export const baseUrl = 'http://localhost:8080/api/v1';

export const useAuthStore = create<IUseAuthStore>(set => ({
	email: '',
	setEmail: (email: string) => set({ email }),
	login: localStorage.getItem('login') || '',
	setLogin: (login: string) => set({ login }),
	password: '',
	setPassword: (password: string) => set({ password }),
	isAuth: JSON.parse(localStorage.getItem('isAuth') || 'false'),
	setIsAuth: (isAuth: boolean) => {
		set({ isAuth });
		localStorage.setItem('isAuth', JSON.stringify(isAuth));
	},
	currentUserLogin: '',
	setCurrentUserLogin: (login: string) => set({ currentUserLogin: login }),

	signIn: async (e, email, password, navigate) => {
		e.preventDefault();

		try {
			const response = await axios.post(`${baseUrl}/auth/signin`, {
				email,
				password,
			});

			console.log('LOG:', response);

			if (response.status === 200) {
				set({ isAuth: true });
				set({ login: response.data.login });

				localStorage.setItem('isAuth', JSON.stringify(true));
				localStorage.setItem('login', response.data.login);

				Cookies.set('session_id', response.data.session_id, {
					expires: 3,
				});

				navigate(`/user/${response.data.login}`);

				console.log('login from lc st: ', localStorage.getItem('login'));

				console.log('response: ', response);
				console.log('login: ', useAuthStore.getState().login);
			}
		} catch (err) {
			console.log('Error signing in', err);
			alert(`Ошибка входа: ${handleError(err)}`);
		}
	},

	signUp: async (e, email, login, password, navigate, signIn) => {
		e.preventDefault();
		try {
			const response = await axios.post(`${baseUrl}/auth/signup`, {
				email,
				login,
				password,
			});

			console.log('up res: ', response);

			if (response.status === 200) {
				set({ login: login });
				localStorage.setItem('login', login);

				signIn(e, email, password, navigate);
			}
		} catch (err) {
			console.log('Error signing up', err);
			alert(`Ошибка регистрации: ${handleError(err)}`);
		}
	},

	logout: async navigate => {
		try {
			const response = await axios.post(
				`${baseUrl}/auth/logout`,
				{},
				{ withCredentials: true }
			);
			if (response.status === 200) {
				set({ isAuth: false, login: '', email: '', password: '' });

				localStorage.setItem('login', '');
				localStorage.setItem('isAuth', JSON.stringify(false));

				Cookies.remove('session_id');

				console.log('after reset login: ', useAuthStore.getState().login);

				navigate('/signin');
			}
		} catch (err) {
			console.log('Error logging out', err);
			alert(`Ошибка выхода: ${handleError(err)}`);
		}
	},
}));
