const fs = require('fs');
const readline = require('readline');

function encodeUrlPartially(url) {
    return url.replace(/ /g, '%20')  // 编码空格
        .replace(/\(/g, '%28')  // 编码左括号
        .replace(/\)/g, '%29'); // 编码右括号
}

async function getData(url, token) {
    const response = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Cookie': `online_token=${token}`
        }
    });

    return response.json();
}

async function fetchData(videoId) {
    const config = JSON.parse(fs.readFileSync('config.json', 'utf8'));
    const token = config.online_token;
    const url = `https://online.njtech.edu.cn/api/v2/videos/${videoId}?id=${videoId}`;

    try {
        const data = await getData(url, token);
        const apiUrl = `https://online.njtech.edu.cn/api/v2/video_seasons/${data.seasons[0]["id"]}?order=asc&orderBy=index`;

        const seasonData = await getData(apiUrl, token);
        const urls = seasonData.episodes.map(episode => episode.url);
        const encodedUrls = urls.map(encodeUrlPartially);

        fs.writeFileSync('urls.txt', encodedUrls.join('\n'), 'utf8');
    } catch (error) {
        console.error(error);
    }
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

rl.question('请输入视频 ID: ', (videoId) => {
    if (!videoId) {
        console.error('视频 ID 不能为空');
        rl.close();
        return;
    }
    fetchData(videoId).finally(() => rl.close());
});
