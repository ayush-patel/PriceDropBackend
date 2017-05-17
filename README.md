# PriceDropBackend
A repository to host Heruko backend code to support the Price Drop iOS App at https://price-drop.herokuapp.com

----

REST API details:
:8080/api/items  | POST  |  send the URL (in JSON) to fetch the details and save as a new product 
:8080/api/items  | GET  |  get the details of all the items saved  

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
"title": "Name 1",
"description": "Brand 1",
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor",
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
"title": "Name 1",
"description": "Brand 1",
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor",
"id": 1,
"originalprice": 195,
"currentprice": 195
},
{
"title": "Name 2",
"description": "Brand 2",
"url": "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor",
"id": 2,
"originalprice": 195,
"currentprice": 195
}
]
