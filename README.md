# rpcstest
RPC stress tester

all test used ./st -url http://{DEV  parity eu_west_1_kitty_flowers ip}/ -n 1000000 -r 100  

1. regualr configuration, about 790 r/s  

limit_req_zone $rate_limit_by zone=limit:10m rate=240r/m;  
limit_req_status 429;  
limit_req zone=limit burst=60;  

MEM: 40% -> 45%  
CPU: 2% -> 20%  
q50: 273ms -> 5ms (due to massive amount of 429)  
q90: 3ms =  
HTTP requests:  
Count: 1.5 -> 800  
200: 1.5 -> 0  
429: 0 -> 800  

2. no limits configuration, about 980 r/s  

limit_req_zone $rate_limit_by zone=limit:10m rate=1000240r/m;  
limit_req_status 429;  
limit_req zone=limit burst=60;  

MEM: 45% -> 47%  
CPU: 5% -> 35%  
q50: 5ms =  
q90: 3ms =  
HTTP requests:  
Count: 2 -> 987  
200: 2 -> 987  
429: 0 =  

3. no_dealy configuration 1, about 980 r/s  

limit_req_zone $rate_limit_by zone=limit:10m rate=58900r/m;  
limit_req_status 429;  
limit_req zone=limit burst=60 nodelay;  

MEM: 44% -> 47%  
CPU: 5% -> 38%  
q50: 5ms =  
q90: 3ms =  
HTTP requests:  
Count: 2 -> 980  
200: 2 -> 980  
429: 0 =  

4. no_dealy configuration 2, about 980 r/s  

limit_req_zone $rate_limit_by zone=limit:10m rate=32000r/m;  
limit_req_status 429;  
limit_req zone=limit burst=60 nodelay;  

MEM: 44% -> 46%  
CPU: 5% -> 30%  
q50: 5ms  
q90: 3ms  
HTTP requests:  
Count: 2 -> 989  
200: 2 -> 537  
429: 0 -> 452  

5. no_dealy configuration 3, about 980 r/s  

limit_req_zone $rate_limit_by zone=limit:10m rate=32000r/m;  
limit_req_status 429;  
limit_req zone=limit burst=3600 nodelay;  

MEM: 44% -> 47%  
CPU: 5% -> 39%  
q50: 5ms =  
q90: 3ms =  
HTTP requests:  
Count: 2 -> 985  
200: 2 -> 537  
429: 0 -> 448  

6. no_dealy configuration 4, about 980 r/s  

limit_req_zone $rate_limit_by zone=limit:10m rate=32000r/m;  
limit_req_status 429;  
limit_req zone=limit burst=13600 nodelay;  

MEM: 44% -> 47%  
CPU: 5% -> 30%  
q50: 5ms  
q90: 3ms  
HTTP requests:  
Count: 2 -> 948  
200: 2 -> 513  
429: 0 -> 435  