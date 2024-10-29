import { useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import { TextField } from '../../components/TextField/TextField';
import { Layout } from '../../layout/Layout';
import { useUsersStore } from '../../store/users/useUsersStore';
import styles from './ChangePasswordPage.module.scss';

export const ChangePasswordPage = () => {
	const login = useUsersStore(state => state.user.login);

	const password = useUsersStore(state => state.password);
	const setPassword = useUsersStore(state => state.setPassword);

	const editPassword = useUsersStore(state => state.editPassword);

	const navigate = useNavigate();

	return (
		<Layout>
			<div className={styles.page}>
				<form>
					<TextField
						className={styles.passwordField}
						direction='row'
						field='input'
						label='Изменить пароль:'
						value={password}
						onChange={e => setPassword(e.target.value)}
						placeholder='Введите новый пароль'
						type='password'
					/>
					<Button
						className={styles.saveButton}
						disabled={!password}
						title={!password ? 'Пароль не может быть пуст' : ''}
						onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
							e.preventDefault();
							editPassword(login, password, navigate);
						}}
					>
						Применить
					</Button>
				</form>
			</div>
		</Layout>
	);
};
