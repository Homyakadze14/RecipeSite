export const deleteLastChar = (str: string) => {
	if (str.slice(-1) == ';') {
		return str.slice(0, -1);
	} else {
		return str;
	}
};

export const addNewLines = (str: string) => {
	if (!str) return null;

	return str.split('\n').map((line, index) => (
		<p key={index}>
			<span>{line}</span>
		</p>
	));
};
