import { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { UserIcon } from '../../assets/icons/UserIcon';
import { Button } from '../../components/Button/Button';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useUsersStore } from '../../store/users/useUsersStore';
import styles from './Header.module.scss';

const AuthButton = () => {
	const logout = useAuthStore(state => state.logout);
	const navigate = useNavigate();

	return (
		<Button
			style={{ width: 106 }}
			className={styles.red}
			onClick={() => logout(navigate)}
		>
			Выйти
		</Button>
	);
};

const LoginLink = () => {
	return (
		<Link to='/signin' style={{ fontSize: 30, fontWeight: 700 }}>
			Войти
		</Link>
	);
};

const UserLink = () => {
	const login = useAuthStore(state => state.login);
	const setParamsLogin = useUsersStore(state => state.setParamsLogin);

	return (
		<Link
			className={styles.userIcon}
			to={`/user/${login}`}
			onClick={() => setParamsLogin(login)}
		>
			<UserIcon />
		</Link>
	);
};

export const Header = () => {
	const login = useAuthStore(state => state.login);
	const isAuth = useAuthStore(state => state.isAuth);
	const paramsLogin = useUsersStore(state => state.paramsLogin);
	const location = useLocation();

	const [path, setPath] = useState(location.pathname);
	const [isBurgerMenuOpen, setIsBurgerMenuOpen] = useState(false);
	const [isScrolled, setIsScrolled] = useState(true);
	const [lastScrollY, setLastScrollY] = useState(0);

	useEffect(() => {
		if (path.includes('/user/')) {
			setPath(`/user/${login}`);
		}
	}, []);

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

	return (
		<header className={`${styles.header} ${!isScrolled && styles.invisible}`}>
			<div>
				<nav className={styles.nav}>
					<ul>
						<li style={{ width: 106 }} />
						<li>
							<Link to='/'>
								<h2 className={styles.title}>Вкусные идеи</h2>
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
									<AuthButton />
								) : (
									<UserLink />
								)
							) : (
								<LoginLink />
							)}
						</li>
					</ul>
				</nav>

				<div className={styles.burgerMenu}>
					<div className={styles.burgerMenuContent}>
						<Link to='/'>
							<h2
								className={styles.title}
								style={{ display: isBurgerMenuOpen ? 'none' : 'block' }}
							>
								Вкусные идеи
							</h2>
						</Link>

						<nav
							className={styles.burgerMenuNav}
							style={{ display: isBurgerMenuOpen ? 'block' : 'none' }}
						>
							{isAuth ? (
								paramsLogin === login && path.includes('/user') ? (
									<AuthButton />
								) : (
									<UserLink />
								)
							) : (
								<LoginLink />
							)}
						</nav>
					</div>
				</div>

				<div
					className={[
						styles.burgerIcon,
						isBurgerMenuOpen && styles.active,
					].join(' ')}
					onClick={() => setIsBurgerMenuOpen(prev => !prev)}
				>
					<span />
				</div>
			</div>
		</header>
	);
};
