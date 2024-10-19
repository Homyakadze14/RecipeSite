import { ILayout } from '../types/interfaces';
import { Header } from './Header/Header';
import styles from './Layout.module.scss';

export const Layout = ({ children }: ILayout) => {
	return (
		<div className={styles.wrapper}>
			<Header />
			<main>{children}</main>
		</div>
	);
};
