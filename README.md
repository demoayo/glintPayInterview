## Top spenders

Output 
- Run command below computes top 5 spenders
- A list of the top 5 spenders and a relative file path to a JSON file containing the top 5 spenders is printed to the terminal
- command args in a JSON Format
    - file_name: specifies path to csv file
    - topN: specifies how many item is included in the ordered topN subset e.g 
        - "10" returns top10 
        - "5" returns top5 
    - filters: specifies which spenders are to be consider in the computation e.g 
        - {"field": "month", "cmp": "=", "value": "2"} Includes spenders where month = 2
        - {"field": "description", "cmp": "=", "value": "CARD SPEND"}  Includes spenders Where descriptions = "CARD SPEND"

```sh
go run main.go '{"file_name": "sample-transactions.csv", "filters": [{"field": "description", "cmp": "=", "value": "CARD SPEND"}, {"field": "month", "cmp": "=", "value": "2"}], "topN": 5}' 

./glintPay '{"file_name": "sample-transactions.csv", "filters": [{"field": "description", "cmp": "=", "value": "CARD SPEND"}, {"field": "month", "cmp": "=", "value": "2"}], "topN": 20}'

```