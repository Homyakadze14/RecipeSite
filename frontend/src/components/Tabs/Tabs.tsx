import styles from './Tabs.module.scss';

export interface ITabs {
	tabs: { [key: string]: string };
	selectedTab: string;
	setSelectedTab: (tab: string) => void;
}

export const Tabs = ({ tabs, selectedTab, setSelectedTab }: ITabs) => {
	return (
		<div className={styles.tabs}>
			<ul>
				{Object.keys(tabs).map(key => (
					<li
						key={key}
						onClick={() => setSelectedTab(key)}
						className={key + ' ' + (key === selectedTab ? styles.active : '')}
					>
						{tabs[key]}
					</li>
				))}
			</ul>
		</div>
	);
};
