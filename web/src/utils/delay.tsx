const delay = (ms: number): Promise<boolean> => new Promise(resolve => setTimeout(() => resolve(true), ms));

export default delay;
