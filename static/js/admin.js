// SPA-like navigation (AJAX)
document.querySelectorAll('.nav-link').forEach(link => {
  link.addEventListener('click', async (e) => {
    e.preventDefault();
    const url = e.target.getAttribute('href');
    const res = await fetch(url);
    const html = await res.text();
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    const newContent = doc.querySelector('#content').innerHTML;
    document.querySelector('#content').innerHTML = newContent;
    history.pushState({}, '', url);
  });
});