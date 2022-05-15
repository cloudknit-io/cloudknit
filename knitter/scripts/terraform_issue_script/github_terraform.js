import fetch from "node-fetch";
import convertor from "json2csv";
import fs from "fs";

const token = "add_your_github_token_here";

export async function getUsers() {

  // Fetch user urls from issue repo

  const userUrls = [];
  try {
    let page = 1;
    while (true) {
      const urls = await urlFetcher(page++);
      const cl = userUrls.length;
      userUrls.push(...new Set(...[urls]).values());
      console.log("pushed " + (userUrls.length - cl));
      if (urls.length === 0) {
        break;
      }
    }
  } catch (ex) {
    console.log(ex);
    fs.writeFileSync("github_user_urls.txt", userUrls.join(";").toString());
  }


  // Use User URLS Fetched, to fetch User Data

  fs.writeFileSync("github_user_urls.txt", userUrls.join(";").toString());
  const urls = fs
    .readFileSync("github_user_urls.txt", {
      encoding: "utf-8",
    })
    .toString()
    .split(";");
  console.log(urls);
  const userData = await Promise.all(urls.map((url) => userFetcher(url)));
  const csv = await convertor.parseAsync(userData);
  fs.writeFileSync("github_terraform_all.csv", csv);
}

async function urlFetcher(page) {
  const data = await fetch(
    `https://api.github.com/repos/hashicorp/terraform/issues?per_page=100&page=${page}`,
    {
      method: "GET",
      headers: [
        ["Accept", "application/vnd.github.v3+json"],
        ["Authorization", `token ${token}`],
      ],
    }
  ).then((res) => res.json());
  return data.map((e) => e.user.url);
}

async function userFetcher(url) {
  return await fetch(url, {
    method: "GET",
    headers: [
      ["Accept", "application/vnd.github.v3+json"],
      ["Authorization", `token ${token}`],
    ],
  }).then((user) => {
    console.log("fetched " + url);
    return user.json();
  });
}
