
// LIST PRODUCTS

function ListProducts(without) {
  fetch("/products/listbyfilter?sale=true", {
    method: "GET",
  }).then((res) => {
    if (res.ok) {
      res = res.json();
      res.then((data) => {
        CreateProductNodes(data, without).then(table => {
          ConfigurePages(table);
          ConfigureTable();
        })
      })
    } else {
      alert("Неудалось получить список продуктов");
    }
  })
}


async function CreateProductNodes(data, without) {
  // Проверка на null или undefined
  if (!data || data.length === 0) {
    return [];
  }
  data = data.filter(pr => {return !without.includes(pr.id)})

  // Сортировка data по полю position
  data.sort((a, b) => a.position - b.position);

  let table = new Array(data.length).fill(undefined); // Создаем массив с нужными индексами, заполненный undefined

  const promises = data.map(async (product, index) => {
    try {
      const photourl = URL.createObjectURL(new Blob([new Int8Array(product.photo)], { type: "image/png" }));


      const newProduct = document.createElement('div');
      newProduct.className = `items-li items-li-${product.id}`;
      newProduct.innerHTML = `
          <p class="items-li-elem items-position">${product.position}</p>
          <div class="items-li-elem items-photo">
            <img src="${photourl}" alt="" class="items-li-elem items-photo-img">
          </div>
          <p class="items-li-elem items-name">${product.name}</p>
          <p class="items-li-elem items-price">${Number(product.price).toFixed(2)}</p>
          <p class="items-li-elem items-discount">${product.sale}%</p>
          <img src="../imgs/icons/make-sale-icon.png" alt="" class="items-li-make-sale" onclick="MakeSale(${product.id}, '${product.name}')">
        `;
      table[index] = newProduct;
    } catch (error) {
      console.error(`Ошибка при обработке товара с ID ${product.id}:`, error);
    }
  });

  await Promise.all(promises);
  return table.filter(item => item !== undefined);
}



function MakeSale(id) {
  let postion=document.querySelector(".items-table-sales").children.length;
  fetch("/admin/sales/add", {
    method: "POST",
    body: JSON.stringify({productId: String(id), position: postion+1})
  }).then(() => ListSales())
}




// LIST SALES
ListSales()

function ListSales() {
  fetch("/sales/list", {
    method: "GET",
  }).then((res) => {
    if (res.ok) {
      return res.json();
    } else {
      alert("Неудалось получить список продуктов");
    }
  }).then(data => {
    localStorage.setItem("salesQuantity", data.length);
    return CreateSalesNodes(data);
  }).then(without => ListProducts(without))
}


function CreateSalesNodes(data) {
  let listWithoutID = []
  let salesCont = document.querySelector(".items-table-sales");
  salesCont.innerHTML = "";
  data.sort((a, b) => a.position - b.position);
  data.forEach(product => {
    listWithoutID.push(product.id)
    const photourl = URL.createObjectURL(new Blob([new Int8Array(product.photo)], { type: "image/png" }));
    const newProduct = document.createElement('div');
    newProduct.className = `items-li items-li-${product.id}`;
    newProduct.innerHTML = `
        <p class="items-li-elem items-position">${product.position}</p>
        <div class="items-li-elem items-photo">
          <img src="${photourl}" alt="" class="items-li-elem items-photo-img">
        </div>
        <p class="items-li-elem items-name">${product.name}</p>
        <p class="items-li-elem items-price">${Number(product.price).toFixed(2)}</p>
        <p class="items-li-elem items-discount">${product.sale}%</p>
        <img src="../imgs/icons/arrow-icon.png" alt="" class="items-li-pos-up" onclick="EditSale(${product.id}, ${product.position-1})">
        <img src="../imgs/icons/arrow-icon.png" alt="" class="items-li-pos-down" onclick="EditSale(${product.id}, ${product.position+1})">
        <img src="../imgs/icons/admin-delete-icon.png" alt="" class="items-li-delete" onclick="DeleteSale(${product.id})">
      `;
    salesCont.appendChild(newProduct);
  });
  return listWithoutID
}


function DeleteSale(id){
  fetch("/admin/sales/delete?id="+id, {
    method: "DELETE"
  }).then(()=>ListSales())
}

function EditSale(id, pos){
  fetch("/admin/sales/edit",{
    method: "PUT",
    body: JSON.stringify({productId:String(id), position: pos})
  }).then(()=>ListSales())

}