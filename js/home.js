
function nextEpisode(id) {
    const url = `update/?id=${id}`
    const response = await fetch(url, { method: "POST" })
    textElement = document.getElementById(`ep${id}`)
    console.log(textElement.innerText)
    episodeString = textElement.innerText
    console.log(episodeString + "1")
    textElement.innerText = episodeString + "1"
}
