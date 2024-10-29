import { IArrowIcon } from '../../types/interfaces';

export const ArrowIcon = ({ direction, fill = '#FFFFFF' }: IArrowIcon) => {
	return (
		<>
			{direction === 'right' ? (
				<svg
					width='20.250000'
					height='21.960938'
					viewBox='0 0 20.25 21.9609'
					fill='none'
					xmlns='http://www.w3.org/2000/svg'
					xmlnsXlink='http://www.w3.org/1999/xlink'
				>
					<defs />
					<path
						id='&gt;'
						d='M0 5.6L14.48 10.92L0 16.35L0 21.96L20.25 13.24L20.25 8.71L0 3.05e-5L0 5.6Z'
						fill={fill}
						fillOpacity='1.000000'
						fillRule='evenodd'
					/>
				</svg>
			) : (
				<svg
					width='19.617188'
					height='21.984375'
					viewBox='0 0 19.6172 21.9844'
					fill='none'
					xmlns='http://www.w3.org/2000/svg'
					xmlnsXlink='http://www.w3.org/1999/xlink'
				>
					<defs />
					<path
						id='&lt;'
						d='M19.61 16.28L5.78 11.03L19.61 5.69L19.61 3.05e-5L0 8.71L0 13.24L19.61 21.98L19.61 16.28Z'
						fill={fill}
						fillOpacity='1.000000'
						fillRule='evenodd'
					/>
				</svg>
			)}
		</>
	);
};
