
export const getSVGNode = (
	attributes: { [key: string]: string },
	svgElementName: string,
	jsProps?: { [key: string]: any }
) => {
	const svgNode: any = document.createElementNS('http://www.w3.org/2000/svg', svgElementName);
	for (const key in attributes) {
		const value = attributes[key];
		svgNode.setAttribute(key, value);
	}

	for (const key in jsProps) {
		if (key) {
			const value = jsProps[key];
			svgNode[key] = value;
		}
	}
	return svgNode;
};
