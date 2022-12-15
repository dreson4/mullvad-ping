# mullvad-ping

Simple script to ping all the mullvad servers written in Golang.

Usage
```
Downloading server list...
There are 897 servers
Albania 1               Australia 25            Austria 15              
Belgium 15              Brazil 4                Bulgaria 6              
Canada 38               Czech Republic 17               Denmark 7               
Estonia 3               Finland 17              France 19               
Germany 44              Greece 1                Hong Kong 14            
Hungary 5               Ireland 4               Israel 3                
Italy 11                Japan 20                Latvia 1                
Luxembourg 4            Moldova 2               Netherlands 20          
New Zealand 3           North Macedonia 1               Norway 20               
Poland 18               Portugal 2              Romania 14              
Serbia 4                Singapore 11            Slovakia 2              
South Africa 2          Spain 16                Sweden 61               
Switzerland 39          UK 60           USA 347                 
United Arab Emirates 1          

Enter 1 to ping all countries.
Enter 2 to ping one country
Action: 2
Country name: Poland

Pinging 185.244.214.58 in Warsaw,Poland..
Pinging 5.253.206.210 in Warsaw,Poland..
Pinging 37.120.156.242 in Warsaw,Poland..
Pinging 37.120.156.162 in Warsaw,Poland..
.
.
.

RESULTS

1 pl-waw-ovpn-202 146.70.144.98 Warsaw Poland -> 282.536ms
2 pl3-wireguard 5.253.206.210 Warsaw Poland -> 283.195667ms
3 pl-waw-wg-201 45.128.38.226 Warsaw Poland -> 283.406ms
4 pl-waw-ovpn-201 146.70.144.66 Warsaw Poland -> 285.104ms
.
.
.

```
