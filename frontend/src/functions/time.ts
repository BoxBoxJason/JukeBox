
const DAYS_OF_WEEK = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
const MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];

export function formatDisplayTime(date: Date | string): string {
  // If the date is a string, convert it to a Date object
  if (typeof date === "string") {
    date = new Date(date);
  }

  const dayOfWeek = DAYS_OF_WEEK[date.getUTCDay()].slice(0, 3);

  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");

  return `${dayOfWeek} ${hours}:${minutes}`;
}

export function fullFormatDisplayTime(date: Date | string): string {
  // If the date is a string, convert it to a Date object
  if (typeof date === "string") {
    date = new Date(date);
  }

  const month = MONTHS[date.getUTCMonth()];
  const day = date.getUTCDate();
  const year = date.getUTCFullYear();

  // Format the hours and minutes
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");
  const seconds = date.getSeconds().toString().padStart(2, "0");

  return `${month} ${day}, ${year} ${hours}:${minutes}:${seconds}`;
}