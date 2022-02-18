make mainnet WITH_ROCKSDB=true
rm -rf /tmp/scf_exchain_data/
exchaind init tmp --chain-id exchain-66 --home /tmp/scf_exchain_data
rm -rf /tmp/scf_exchain_data/data
cp -rf /Users/oker/scf/data/s0-5810700-rocksdb/data /tmp/scf_exchain_data
exchaind replay -d /Users/oker/scf/data/sx-5811000-5813000-rocksdb/data --home /tmp/scf_exchain_data --paralleled-sender=true --multi-cache=true  > scf.log 2>&1 &
