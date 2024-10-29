import { IStarIcon } from '../../types/interfaces';

export const StarIcon = ({
	fill = '#DC780B',
	width = 16,
	height = 15.4,
	onClick,
}: IStarIcon) => {
	return (
		<svg
			width={width}
			height={height}
			viewBox='0 0 16 15.4019'
			fill='none'
			xmlns='http://www.w3.org/2000/svg'
			xmlnsXlink='http://www.w3.org/1999/xlink'
			onClick={onClick}
		>
			<defs />
			<path
				id='Vector'
				d='M16 5.88L10.18 5.5L7.99 0L5.8 5.5L0 5.88L4.45 9.66L2.99 15.4L7.99 12.23L13 15.4L11.53 9.66L16 5.88Z'
				fill={fill}
				fillOpacity='1.000000'
				fillRule='nonzero'
			/>
		</svg>
	);
};
