import { useState } from 'react';
import CopyToClipboard from 'react-copy-to-clipboard';
import { Button } from '../../components/Button/Button';
import { Layout } from '../../layout/Layout';
import { useTokenStore } from '../../store/token/useTokenStore';
import styles from './GenTokenPage.module.scss';

export const GenTokenPage = () => {
	const token = useTokenStore(state => state.token);
	const getToken = useTokenStore(state => state.getToken);

	const [isCopied, setIsCopied] = useState(false);

	const handleCopy = () => {
		setIsCopied(true);
		navigator.clipboard.readText().then(copiedText => {
			if (copiedText === token) {
				console.log('Text was successfully copied!');
			} else {
				console.log('Text was not copied correctly.');
			}
		});
	};

	return (
		<Layout>
			<div className={styles.page}>
				<span className={styles.title}>Ваш токен: </span>

				<div className={styles.tokenContainer}>
					{token ? (
						<>
							<span className={styles.token}>{token}</span>
							<CopyToClipboard text={token} onCopy={handleCopy}>
								<Button className={styles.copyButton}>
									<svg
										style={{ height: 40, width: 40 }}
										xmlns='http://www.w3.org/2000/svg'
										width='24'
										height='24'
										viewBox='0 0 24 24'
										fill='none'
										stroke='currentColor'
										strokeWidth='2'
										strokeLinecap='round'
										strokeLinejoin='round'
									>
										<rect
											x='9'
											y='9'
											width='13'
											height='13'
											rx='2'
											fill='none'
										/>
										<path d='M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h3' />
										<path d='M15 5v2a2 2 0 0 0 2 2h3' />
									</svg>
								</Button>
							</CopyToClipboard>
						</>
					) : (
						<span className={styles.empty}>Вы ещё не сгенерировали токен</span>
					)}
				</div>
				<span className={styles.copied}>{isCopied ? 'Скопировано!' : ''}</span>
				<Button
					className={styles.genButton}
					disabled={Boolean(token)}
					title={token ? 'Токен уже сгенерирован' : ''}
					onClick={getToken}
				>
					Сгенерировать
				</Button>
			</div>
		</Layout>
	);
};
