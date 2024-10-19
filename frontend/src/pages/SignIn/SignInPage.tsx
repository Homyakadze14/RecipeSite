import { Link, useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import { Input } from '../../components/Input/Input';
import { useAuthStore } from '../../store/auth/useAuthStore';
import styles from './SignInPage.module.scss';

export const SignInPage = () => {
	const email = useAuthStore(state => state.email);
	const setEmail = useAuthStore(state => state.setEmail);

	const password = useAuthStore(state => state.password);
	const setPassword = useAuthStore(state => state.setPassword);

	const signIn = useAuthStore(state => state.signIn);

	const navigate = useNavigate();

	return (
		<div className={styles.wrapper}>
			<main className={styles.main}>
				<div className={styles.container}>
					<h2>Вкусные идеи</h2>
					<div>
						<h3>Ещё нет аккаунта?</h3>
						<Link to='/signup'>Зарегистрироваться</Link>
					</div>
					<form>
						<Input
							value={email}
							onChange={e => setEmail(e.target.value)}
							placeholder='Email'
						/>
						<Input
							value={password}
							onChange={e => setPassword(e.target.value)}
							placeholder='Password'
							type='password'
						/>
						<Button
							disabled={Object.values({ email, password }).some(
								value => !value
							)}
							title={
								Object.values({ email, password }).some(value => !value)
									? 'Все поля должны быть заполнены'
									: ''
							}
							onClick={(e: React.MouseEvent<HTMLButtonElement>) =>
								signIn(e, email, password, navigate)
							}
						>
							Войти
						</Button>
					</form>
				</div>
			</main>
		</div>
	);
};
