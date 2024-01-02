Questions for seQura

1. In the "Problem Statement" the instructions state that the new disbursement payout system for seQura commissions must be calculated for "existing, present in the CSV files, and new orders". Does this mean that existing orders will always be in the CSV files, and that new orders may arrive and be inserted into the existing CSV files, or that perhaps there may be updated files containing new orders that will be expected to be ingested into the new disbursement payout system?

2. At what time in UTC will the arrival of any new orders be rolled to the next date? The requirements state that the system must process all disbursements by 8:00 UTC. I am wondering if the CSV file provided is an accurate representation of the number of records the production system will need to process, or only a portion of the data. If it is only a portion, can you tell me what percentage of the data the sample represents?

3. Reporting and record look up requirements aren't stated, but the "Problem Statement" does state that "Disbursements groups all the orders for a merchant in a given day or week." . After reviewing the orders.csv file it appears that the orders are entered into the file and grouped by merchant. Will the orders.csv file always be ordered by merchant, or can the file contain time/date based entries where the merchant is random?

4. Will the dispersement system be used as a source for reporting directly on its data, or will there be an upstream system ingesting its output to later run quires against?