type Args = string | number;

// eslint-disable-next-line no-undef
const env = process.env.NODE_ENV;

export class TemplateHelper {
	private static replace(regex: RegExp, templateName: string, ...args: Args[]): string {
		const _arguments = Array.prototype.slice.call(args);
		const len = _arguments.length;
		if (!templateName || len === 0) {
			return templateName;
		}

		let i = 0;
		return templateName.replace(regex, function () {
			const value = _arguments[i];
			if (env === 'development' && !value && value !== '') {
				throw new Error(`Missing argument for TEMPLATE: ${templateName}`);
			}
			i++;
			return value;
		});
	}

	public static format(templateName: string, ...args: Args[]): string {
		return TemplateHelper.replace(/%\w+%/g, templateName, ...args);
	}

	public static route(templateName: string, ...args: Args[]): string {
		return TemplateHelper.replace(/:\w+/g, templateName, ...args);
	}
}
