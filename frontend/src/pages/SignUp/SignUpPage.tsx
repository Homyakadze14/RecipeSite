import { Link, useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import { Input } from '../../components/Input/Input';
import { useAuthStore } from '../../store/auth/useAuthStore';
import styles from './SignUpPage.module.scss';

export const SignUpPage = () => {
	const email = useAuthStore(state => state.email);
	const setEmail = useAuthStore(state => state.setEmail);

	const login = useAuthStore(state => state.login);
	const setLogin = useAuthStore(state => state.setLogin);

	const password = useAuthStore(state => state.password);
	const setPassword = useAuthStore(state => state.setPassword);

	const signUp = useAuthStore(state => state.signUp);
	const signIn = useAuthStore(state => state.signIn);

	const navigate = useNavigate();

	return (
		<div className={styles.wrapper}>
			<main className={styles.main}>
				<div className={styles.container}>
					<h2>Вкусные идеи</h2>
					<div>
						<h3>Уже есть аккаунт?</h3>
						<Link to='/signin'>Войти</Link>
					</div>
					<form>
						<Input
							value={email}
							onChange={e => setEmail(e.target.value)}
							placeholder='Email'
						/>
						<Input
							value={login}
							onChange={e => setLogin(e.target.value)}
							placeholder='Login'
						/>
						<Input
							value={password}
							onChange={e => setPassword(e.target.value)}
							placeholder='Password'
							type='password'
						/>
						<Button
							disabled={Object.values({ email, login, password }).some(
								value => !value
							)}
							title={
								Object.values({ email, login, password }).some(value => !value)
									? 'Все поля должны быть заполнены'
									: ''
							}
							onClick={(e: React.MouseEvent<HTMLButtonElement>) =>
								signUp(e, email, login, password, navigate, signIn)
							}
						>
							Зарегистрироваться
						</Button>
					</form>
				</div>
			</main>
		</div>
	);
};
