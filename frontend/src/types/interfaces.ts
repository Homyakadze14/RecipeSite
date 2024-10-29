import { NavigateFunction } from 'react-router-dom';
import { Filter } from './types';

export interface IRecipe {
	id: number;
	about: string;
	complexity: 1 | 2 | 3;
	created_at: string;
	ingridients: string;
	instructions: string;
	need_time: string;
	title: string;
	photos_urls: string;
	updated_at: string;
	author: IAuthor;
}

export interface IAuthor {
	login: string;
	icon_url: string;
}

export interface IUser {
	id: number;
	login: string;
	about: string;
	icon_url: string;
	created_at: string;
	recipies: IRecipe[];
	liked_recipies: IRecipe[];
	is_subscribed: boolean;
}

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

export interface IArrowIcon extends React.SVGProps<SVGSVGElement> {
	direction: 'left' | 'right';
	fill?: string;
}

export interface ILikeIcon extends React.SVGProps<SVGSVGElement> {
	fill?: string;
}

export interface IStarIcon extends React.SVGProps<SVGSVGElement> {
	fill?: string;
}

export interface IComment {
	id: number;
	author: IAuthor;
	text: string;
}

export interface ICommentsList {
	comments: IComment[];
}

export interface IImageUploader {
	onImageUpload: (file: File | null, imageSrc: string | null) => void;
	selectedImage: File | null;
	label: string;
}

export interface IInput {
	value: string;
	onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
	placeholder?: string;
	type?: string;
}

export interface IModal extends React.HTMLAttributes<HTMLDivElement> {
	isActive: boolean;
	setIsActive: React.Dispatch<React.SetStateAction<boolean>>;
	children: React.ReactNode;
}

export interface IRecipesList {
	data: IRecipe[];
	noDataText?: string;
	isSearchQueryChanged?: boolean;
	isTabChanged?: boolean;
	listStartY?: number;
}

export interface ISearchInput {
	value: string;
	onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export interface ITabs {
	tabs: { [key: string]: string };
	selectedTab: string;
	setSelectedTab: (tab: string) => void;
}

export interface ITextField {
	direction: 'row' | 'column';
	label: string;
	field: 'input' | 'textarea';
	value: string;
	onChange: (
		event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => void;
	placeholder?: string;
	type?: string;
	className?: string;
}

export interface IUserCard {
	avatarSrc: string;
	name: string;
	description: string;
}

export interface IEditRecipePage {
	recipeId: number;
}

export interface IUseAuthStore {
	email: string;
	setEmail: (email: string) => void;
	login: string;
	setLogin: (login: string) => void;
	password: string;
	setPassword: (password: string) => void;
	isAuth: boolean;
	setIsAuth: (isAuth: boolean) => void;
	currentUserLogin: string;
	setCurrentUserLogin: (login: string) => void;

	signUp: (
		e: React.MouseEvent<HTMLButtonElement>,
		email: string,
		login: string,
		password: string,
		navigate: (to: string) => void,
		signIn: (
			e: React.MouseEvent<HTMLButtonElement>,
			email: string,
			password: string,
			navigate: (to: string) => void
		) => void
	) => void;
	signIn: (
		e: React.MouseEvent<HTMLButtonElement>,
		email: string,
		password: string,
		navigate: (to: string) => void
	) => void;
	logout: (navigate: NavigateFunction) => void;
}

export interface IUseRecipesStore {
	isLoading: boolean;
	setIsLoading: (isLoading: boolean) => void;
	isLiked: boolean;
	setIsLiked: (isLiked: boolean) => void;
	likes_count: number;
	setLikesCount: (likes_count: number) => void;
	recipe: IRecipe;
	setRecipe: (recipe: IRecipe) => void;
	recipes: IRecipe[];
	setRecipes: (recipes: IRecipe[]) => void;
	recipeId: number;
	setRecipeId: (recipeId: number) => void;

	comments: IComment[];
	setComments: (comments: IComment[]) => void;

	text: string;
	setText: (text: string) => void;

	createRecipeForm: {
		title: string;
		about: string;
		ingridients: string;
		author: IAuthor;
		instructions: string;
		complexity: string | number;
		need_time: string;
		selectedImage: any;
	};
	setCreateRecipeForm: (form: any) => void;

	editRecipeForm: {
		photos: File | null;
		title: string;
		about: string;
		ingridients: string;
		instructions: string;
		complexity: string | number;
		need_time: string;
	};
	setEditRecipeForm: (form: any) => void;

	createRecipe: (
		e: React.MouseEvent<HTMLButtonElement>,
		login: string,
		createRecipeForm: {
			title: string;
			about: string;
			ingridients: string;
			author: IAuthor;
			instructions: string;
			complexity: string | number;
			need_time: string;
			selectedImage: any;
		},
		navigate: NavigateFunction
	) => void;

	getAllRecipes: () => void;
	getFilteredRecipes: (filter: Filter) => void;
	getRecipeAndComments: (recipeId: number) => Promise<void>;
	getAuthor: (login: string) => Promise<string>;
	editRecipe: (
		e: React.MouseEvent<HTMLButtonElement>,
		login: string,
		recipeId: number,
		editRecipeForm: any,
		navigate: NavigateFunction
	) => void;
	deleteRecipe: (login: string, recipeId: number) => void;

	likeRecipe: (recipeId: number) => void;
	unlikeRecipe: (recipeId: number) => void;

	createComment: (
		e: React.MouseEvent<HTMLButtonElement>,
		recipeId: number,
		text: string
	) => void;
	editComment: (commentId: number, text: string) => void;
	deleteComment: (commentId: number, recipeId: number) => void;
}

export interface IUseTokensStore {
	token: string;
	setToken: (token: string) => void;
	getToken: () => void;
}

export interface IUseUsersStore {
	user: IUser;
	setUser: (user: IUser) => void;
	users: IUser[];
	setUsers: (users: IUser[]) => void;

	paramsLogin: string;
	setParamsLogin: (login: string) => void;

	editUserForm: {
		icon: any;
		login: string;
		about: string;
	};
	setEditUserForm: (form: any) => void;

	password: string;
	setPassword: (password: string) => void;

	getUser: (login: string) => void;

	editUser: (
		e: React.MouseEvent<HTMLButtonElement>,
		login: string,
		editUserForm: { icon: File | null; login: string; about: string },
		navigate: NavigateFunction,
		setIsEditModalVisible: (isEditModalVisible: boolean) => void
	) => void;

	editPassword: (
		login: string,
		newPassword: string,
		navigate: NavigateFunction
	) => void;

	subscribe: (login: string) => void;
	unsubscribe: (login: string) => void;
}
