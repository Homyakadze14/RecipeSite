import { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { UserIcon } from '../../assets/icons/UserIcon';
import { Button } from '../../components/Button/Button';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useUsersStore } from '../../store/users/useUsersStore';
import styles from './Header.module.scss';

export const Header = () => {
	const isAuth = useAuthStore(state => state.isAuth);

	const login = useAuthStore(state => state.login);
	const paramsLogin = useUsersStore(state => state.paramsLogin);

	const logout = useAuthStore(state => state.logout);

	const navigate = useNavigate();

	const location = useLocation();

	const [path, setPath] = useState(location.pathname);

	useEffect(() => {
		if (path.includes('/user/')) {
			setPath(`/user/${login}`);
		}
	}, [login]);

	const setParamsLogin = useUsersStore(state => state.setParamsLogin);

	const [isScrolled, setIsScrolled] = useState(true);
	const [lastScrollY, setLastScrollY] = useState(0);

	useEffect(() => {
		const handleScroll = () => {
			const currentScrollY = window.scrollY;

			if (
				currentScrollY >
				(path.includes('/user') ? 330 : path === '/' ? 100 : 60)
			) {
				if (currentScrollY > lastScrollY) {
					setIsScrolled(false);
				} else {
					setIsScrolled(true);
				}
			} else {
				setIsScrolled(true);
			}

			setLastScrollY(currentScrollY);
		};

		window.addEventListener('scroll', handleScroll);

		return () => {
			window.removeEventListener('scroll', handleScroll);
		};
	}, [lastScrollY]);

	useEffect(() => {
		console.log('LOGIN LOGIN LOGIN LOGIN: ', login);
		console.log('PATH PATH PATH PATH: ', path);
	}, [login, path]);

	return (
		<header className={`${styles.header} ${!isScrolled && styles.invisible}`}>
			<div>
				<nav>
					<ul>
						<li style={{ width: 106 }}></li>
						<li>
							<Link to='/'>
								<h2>Вкусные идеи</h2>
							</Link>
						</li>
						<li
							style={{
								width: 106,
								display: 'flex',
								justifyContent: 'flex-end',
							}}
						>
							{isAuth ? (
								paramsLogin === login && path.includes('/user') ? (
									<Button
										className={styles.red}
										onClick={() => logout(navigate)}
									>
										Выйти
									</Button>
								) : (
									<Link
										to={`/user/${login}`}
										onClick={() => setParamsLogin(login || '')}
									>
										<UserIcon />
									</Link>
								)
							) : (
								<Link to='/signin'>Войти</Link>
							)}
						</li>
					</ul>
				</nav>
			</div>
		</header>
	);
};
