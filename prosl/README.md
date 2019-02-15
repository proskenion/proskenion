# Proskenion Specific Language based on YAML

Prosl means (Proskenion specific language) based on yaml. Data is expressed by protocol buffers. So, protofile is designed this DSL.

[Based on prosl.proto](https://github.com/proskenion/proskenion/blob/master/proto/prosl.proto)

## prosl convertor

Yaml file convert to protobuf format.

```
$ ./proslc prosl.yaml
```

## prosl validator

Yaml file validate(type check).

```
$ ./proslv prosl.yaml
```

## For example to write yaml
### genesis
```yaml
- set:
    - tx
    - transaction:
        commands:
          - add_peer:
              authorizer_id: root@com
              peer_id: root@peer
              address: 127.0.0.1:50055
              public_key: 0x3788ef7f97cbc4bda223add5ea147fa3e8a096ad4f27b0dcf247e9fb9443060e
          - create_account:
              authorizer_id: root@com
              account_id: authorizer@com
              public_keys:
                list: nil
              quorum: 0
          - create_account:
              authoirzer_id: root@com
              account_id: incentive@com
              public_keys:
                list: nil
              quorum: 0
          - define_storage:
              authorizer_id: root@root
              storage_id: /degraders
              storage:
                storage:
                  acs:
                    list: nil
          - create_storage:
              authorizer_id: root@root
              wallet_id: root@root/degraders
          - consign:
              authorizer_id: root@root
              account_id: authorizer@com
              peer_id: root@peer
          - consign:
              authorizer_id: root@root
              account_id: incentive@com
              peer_id: root@peer
          - define_storage:
              authorizer_id: root@root
              storage_id: /prflag
              storage:
                storage:
                  prosl_type: "none"
- return:
    variable: tx
```

### consensus
```yaml
- set:
    - acs
    - query:
        authorizer: root@com
        select: "*"
        type: List
        from: com/account
        order_by:
          - balance
          - DESC
        limit: 20
- return:
    variable: acs
```

### incentive
```yaml
- set:
    - height
    - valued:
        - variable: top
        - int64
        - height
- set: # height % 20 を取得、20回で周期的
    - ind
    - cast:
        - int32
        - mod:
            - variable: height
            - 2ll
- if: # 一周したら
    - eq:
        - variable: ind
        - 0
    - set: # account List を取得
        - acs
        - query:
            authorizer: root@root
            select: "*"
            type: List
            from: pr/account
            order_by:
              - balance
              - DESC
            limit: 4
    - set:
        - add_degrade
        - update_object:
            authorizer_id: root@root
            wallet_id: root@root/degraders
            key: acs
            object:
              variable: acs
    - return:
        transaction:
          commands:
            - variable: add_degrade
            - add_balance:
                authorizer_id: root@root
                account_id:
                  valued:
                    - indexed:
                        - variable: acs
                        - account
                        - 0
                    - address
                    - id
                balance: 10000ll
- else:
    - set:
        - acs
        - query:
            authorizer: root@root
            select: acs
            type: List
            from: root@root/degraders
    - return:
        transaction:
          commands:
            - add_balance:
                authorizer_id: root@root
                target_id:
                  valued:
                    - indexed:
                        - variable: acs
                        - account
                        - 1
                    - address
                    - id
                balance: 10000ll
```
