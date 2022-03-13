## Top spenders

- Run the command below to computes the top 5 spenders.
- Service outputs a list of the top 5 spenders and a relative file path to a JSON file containing the top 5 spenders is printed to the terminal
- command args in a JSON Format
    - file_name: specifies path to csv file
    - top_n: specifies how many item is included in the ordered topN subset e.g 
        - "10" returns top 10 
        - "5" returns top 5 
    - filters: specifies which spenders are to be consider in the computation e.g 
        - {"field": "month", "cmp": "=", "value": "2"} Includes spenders where month = 2
        - {"field": "description", "cmp": "=", "value": "CARD SPEND"}  Includes spenders Where descriptions = "CARD SPEND"

```sh
go run . '{"file_name": "sample-transactions.csv", "filters": [{"field": "description", "cmp": "=", "value": "CARD SPEND"}, {"field": "month", "cmp": "=", "value": "2"}], "top_n": 5}' 

./glintPay '{"file_name": "sample-transactions.csv", "filters": [{"field": "description", "cmp": "=", "value": "CARD SPEND"}, {"field": "month", "cmp": "=", "value": "2"}], "top_n": 20}'

```