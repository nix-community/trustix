export class NameValuePair<T> {
  public name: string;
  public value: T;

  constructor(name: string, value: T) {
    this.name = name;
    this.value = value;
  }

  // Turn a map of objects into an array of name value pairs
  public static fromMap<T>(values: { [key: string]: T }): NameValuePair<T>[] {
    const arr: Array<T> = [];

    Object.keys(values).forEach((key) => {
      const value = values[key];
      arr.push(new NameValuePair(key, value));
    });

    return arr;
  }
}
