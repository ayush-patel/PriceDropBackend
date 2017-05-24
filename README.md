# PriceDropBackend
A repository to host Heruko backend code to support the Price Drop iOS App at https://price-drop.herokuapp.com

----

REST API details:

:PORT/api/items  | POST  |  send the URL (in JSON) to fetch the details and save as a new product 

:PORT/api/items  | GET  |  get the details of all the items saved  

----

Sample POST:
Headers: 
Content-Type:application/json

Body: 
{
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor"
}

Response:
{
"title": "Gemini Link Tote",
"description": "Nordstrom",
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor",
"imageurl": "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/0/_13480640.jpg?crop=pad&pad_color=FFF&format=jpeg&trim=color&trimcolor=FFF&w=60&h=90",
"id": 1,
"originalprice": 195,
"currentprice": 195
}

----

Sample GET:
Headers: Content-Type:application/json

Response:
[
{
"title": "Gemini Link Tote",
"description": "Nordstrom",
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor",
"imageurl": "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/0/_13480640.jpg?crop=pad&pad_color=FFF&format=jpeg&trim=color&trimcolor=FFF&w=60&h=90",
"id": 1,
"originalprice": 195,
"currentprice": 195
},
{
"title": "Gemini Link Tote",
"description": "Nordstrom",
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor",
"imageurl": "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/0/_13480640.jpg?crop=pad&pad_color=FFF&format=jpeg&trim=color&trimcolor=FFF&w=60&h=90",
"id": 2,
"originalprice": 195,
"currentprice": 195
}
]
