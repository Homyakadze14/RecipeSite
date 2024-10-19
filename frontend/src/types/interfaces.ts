import { IAuthor } from '../store/users/useUsersStore';

export interface ILayout {
	children: React.ReactNode;
}

export interface IComplexity {
	starsAmount: string | number;
	setStars?: (stars: number | string) => void;
	width?: number;
	height?: number;
	isClickable?: boolean;
}

export interface IButton extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	className?: string;
	children: React.ReactNode;
	onClick?: React.MouseEventHandler<HTMLButtonElement>;
}

export interface IRecipeCard {
	id: number;
	imageSrc: string;
	title: string;
	description: string;
	author: IAuthor;
	time: string;
	starsAmount: 1 | 2 | 3;
	onClick?: () => void;
	onDelete: (login: string, recipeId: number) => void;
}

export interface IRecipeCardDetail {
	id: number;
	photos_urls: string;
	title: string;
	about: string;
	need_time: string;
	complexity: 1 | 2 | 3;
	author: IAuthor;
	ingredients: string;
	instructions: string;
}
