import { useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import { Layout } from '../../layout/Layout';
import styles from './NotFoundPage.module.scss';

export const NotFoundPage = () => {
	const navigate = useNavigate();

	return (
		<Layout>
			<div className={styles.notFoundPage}>
				<h2>404</h2>
				<p>
					Упс, кажется, вы немного заблудились!
					<br />
					Возвращайтесь на главную страничку скорее, вас там ждут новые рецепты.
				</p>
				<Button className={styles.button} onClick={() => navigate('/')}>
					Вернуться на главную
				</Button>
			</div>
		</Layout>
	);
};
