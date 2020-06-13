/*
/* pagination
*/

// input: totalItems, currentPage, pageSize, maxPages
function paginate(e,t=1,a=10,r=10){let n,l,g=Math.ceil(e/a);if(t<1?t=1:t>g&&(t=g),g<=r)n=1,l=g;else{let e=Math.floor(r/2),a=Math.ceil(r/2)-1;t<=e?(n=1,l=r):t+a>=g?(n=g-r+1,l=g):(n=t-e,l=t+a)}let i=(t-1)*a,s=Math.min(i+a-1,e-1),o=Array.from(Array(l+1-n).keys()).map(e=>n+e);return{totalItems:e,currentPage:t,pageSize:a,totalPages:g,startPage:n,endPage:l,startIndex:i,endIndex:s,pages:o}}

const pageSize = 6;
const maxPages = 3;
const currentPage = new URL(window.location.href).searchParams.get("page");

function pageToOffset(page) {
    const result = (page - 1) * pageSize;
    if (result < 0) return 0;
    return result;
}

function offsetToPage(offset) {
    return (offset / pageSize) + 1;
}

function changePage(page) {
    window.location = window.location.origin
        + window.location.pathname
        + "?offset=" + pageToOffset(page)
        + "&limit=" + pageSize;
}

// change page if page param was changed
// if (currentPage && currentPage != offsetToPage(backEndOffset)) {
//     changePage(currentPage);
// }

// change offset to page on url
// const currentOffset = new URL(window.location.href).searchParams.get("offset");
// if (currentOffset) {
//     history.pushState(null, null, "?page=" + offsetToPage(currentOffset));
// }

const paginationResult = paginate(total, parseInt(currentPage || 1), pageSize, maxPages);
const pageNumbersElement = document.getElementById("page-numbers");
const goToUrl = window.location.origin + window.location.pathname + "?page=";
const currentPageFromOffset = offsetToPage(backEndOffset);

// create buttons
if (paginationResult.totalItems > limit) createPaginationButtons();
else document.getElementById("pagination").className = "hidden";

function createPaginationButtons() {
    paginationResult.pages.forEach(function (value, index) {
        const node = document.createElement("a");

        node.href = "#";
        node.innerText = value;
        node.className = currentPageFromOffset === value ? "active" : "";

        pageNumbersElement.appendChild(node)
    });

    pageNumbersElement.addEventListener("click", function (e) {
        e.preventDefault();
        changePage(e.target.childNodes[0].data);
    })

    // prev
    const prevElement = document.querySelector(".prev");
    prevElement.className += (paginationResult.startPage === currentPageFromOffset) ? " block" : "";

    prevElement.addEventListener("click", function (e) {
        e.preventDefault();
        const prevPage = currentPageFromOffset - 1;
        if (prevPage < paginationResult.startPage) return;

        changePage(prevPage);
    })

    // next
    const nextElement = document.querySelector(".next");
    nextElement.className += (paginationResult.endPage === currentPageFromOffset) ? " block" : "";

    nextElement.addEventListener("click", function (e) {
        e.preventDefault();
        const nextPage = offsetToPage(backEndOffset) + 1;
        if (nextPage > paginationResult.endPage) return;

        changePage(nextPage);
    })
}


/*
/* select
*/
let pathname = window.location.pathname.split("/");
if (pathname.length === 2 && pathname[1] !== "") {
    document.getElementById("select").selectedIndex = document.querySelector('option[value="' + pathname[pathname.length - 1] + '"]').index;
}

// busca
let busca = document.getElementById("busca").addEventListener("keydown", function(e) {
    if (e.keyCode === 13 || e.which === 13) {
        window.location = "/search/" + document.getElementById("busca").value.trim().toLowerCase();
    }
});

function goToSquad() {
    var selectElement = document.getElementsByTagName("select")[0];
    window.location = "/" + document.getElementsByTagName("select")[0].options[selectElement.selectedIndex].value;
}
