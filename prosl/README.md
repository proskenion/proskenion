# Proskenion Specific Language based on YAML

Prosl means (Proskenion specific language) based on yaml.


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

```yaml
- set: # account List を取得
    - account
    - query:
        select: *
        type: List
        from: domain.com
        where:
            gt:
                - reputation
                - 0.5
        order_by:
            - fav
            - DESC

- set: # height % 20 を取得、20回で周期的
    - ind
    - %:
        - value:
            - top
            - int32
            - height
        - 20
- if: # 一周したら
    - eq:
        - ind
        - 0
    - set:
        - degraders
        - query:
            select: *
            type: List
            from: domain.com
            where:
                gt:
                    - reputation
                    - 0.5
            order_by:
                - fav
                - DESC
            - limit: 20
    - set:
        - add_degrade
        - update:
            target_id: root@domain.com#degrader.accounts
            objects: degraders
    - return:
        transaction:
            commands:
                - add_degrade
else:
    - return:
        transaction:
            created_time: nil
```

```yaml
- set:
    - peers
    - query:
        select: peer
        type: Peer
        from: domain.com#degrader.accounts
        order_by:
            - fav
            - DESC
        - limit:20
- return: peers
```
