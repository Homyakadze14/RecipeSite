import { useState } from 'react';
import { StarIcon } from '../../assets/icons/StarIcon';
import { IComplexity } from '../../types/interfaces';
import styles from './Complexity.module.scss';

export const Complexity = ({
	starsAmount = 1,
	setStars,
	width,
	height,
	isClickable = false,
}: IComplexity) => {
	const [stars, setStarsState] = useState<number>(Number(starsAmount));

	const maxStarsAmount = 3;

	const handleClick = (index: number) => {
		console.log('handleClick', index);
		if (isClickable && setStarsState) {
			const newStars = index + 1;

			console.log('newStars: ', newStars);
			setStarsState(newStars);
			console.log('stars: ', stars);

			if (setStars) {
				setStars(newStars);
			}
		}
	};

	return (
		<div className={styles.complexity}>
			{Array.from({
				length: stars,
			}).map((_, index) => (
				<StarIcon
					key={index}
					fill={stars === 1 ? '#46EE0B' : stars === 2 ? '#FF9900' : 'red'}
					onClick={() => handleClick(index)}
					width={width}
					height={height}
					pointerEvents='auto'
				/>
			))}
			{Array.from({ length: maxStarsAmount - stars }).map((_, index) => (
				<StarIcon
					key={index}
					fill='#EDEDED'
					onClick={() => handleClick(index + stars)}
					width={width}
					height={height}
					pointerEvents='auto'
				/>
			))}
		</div>
	);
};
