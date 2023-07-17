# Blockchain-Event-Handling-System
## KCC 2023 학부생논문 - 허가형 블록체인의 오프체인 데이터 형성 방법

> #### 하이퍼레저 패브릭의 실행 환경에서 발생하는 이벤트를 수신하여
> #### 오프체인 데이터를 형성하는 이벤트 핸들링 시스템입니다.
> #### 하이퍼레저 패브릭에 대한 확장성 및 개방성 증대를 목표로 하였습니다.

### 목차
[1. 소개](#주요-기능)<br>
[2. 사용 기술](#사용-기술)<br>
[3. 아키텍처](#아키텍처-및-흐름도)<br>
[4. 개발](#개발)<br>
&emsp;[4-1. 코드](#코드)<br>
&emsp;[4-2. 화면](#)<br>
[5. 성능 평가](#성능-평가)<br>
[6. 결과](#결과)<br><br>

## 주요 기능
허가형 블록체인의 실행환경에서 발생하는 이벤트(Contract Event, Block Event)를 수신하여 시간,공간,참가자,트랜잭션 데이터로 분류하여 오프체인 데이터베이스에 저장한다.
이를 통해 허가형 블록체인에 적용 가능한 분류별 데이터로 외부 서드파티 시스템과 함께 활용할 수 있으며, 블록체인 네트워크에 발생하는 부하를 해소할 수 있다.

## 사용 기술
- 블록체인 : Hyperledger Fabric, Fabric Gateway SDK
- 시스템 : Gin(Golang)
- 오프체인 데이터베이스 : MySQL
- 테스트 : 쉘스크립트, pandas

## 아키텍처 및 흐름도
|아키텍처|시퀀스 다이어그램|
|---|---|
|<img width=400 src=https://github.com/jhl8109/FabricAPI/assets/78259314/4bda1ad2-998e-4cd1-8559-82973c36a9d4/>|![New Architecture (5)](https://github.com/jhl8109/FabricAPI/assets/78259314/b430dcfb-daf2-4944-8a47-bebf3c25166e)|
<br>


## 개발
### 코드
#### 게이트웨이 연결
```golang
func Connect() {
	clientConnection := newGrpcConnection()

	id := newIdentity()
	sign := newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	network = gateway.GetNetwork(channelName)
	contract = network.GetContract(chaincodeName)
	fmt.Printf("*** first:%s\n", contract)

	ctx, _ := context.WithCancel(context.Background())

	startChaincodeEventListening(ctx, network)
}
```
---
#### 이벤트 핸들링
```golang
func startChaincodeEventListening(ctx context.Context, network *client.Network) {

	blockEvents, blockErr := network.BlockEvents(ctx, client.WithStartBlock(1))
	if blockErr != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", blockErr))
	}
	fmt.Println("\n*** Start Block event listening")

	ccEvents, ccErr := network.ChaincodeEvents(ctx, chaincodeName)
	if ccErr != nil {
		panic(fmt.Errorf("failed to start block event listening: %w", ccErr))
	}
	fmt.Println("\n*** Start chaincode event listening")

	go func() {
		for event := range blockEvents {
			hashBytes := event.GetHeader().GetDataHash()
			hashString := fmt.Sprintf("%x", hashBytes)
			blockNumber := event.GetHeader().GetNumber()
			fmt.Printf("\n<-- Block event received: \n   Received block number : %d \n   Received block hash - %s\n", blockNumber, hashString)
		}
	}()
	go func() {
		outputFile := "process.txt"
		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		for event := range ccEvents {
			startTime := time.Now().UnixNano()
			startTimeString := fmt.Sprintf("%d", startTime)

			eventStr := formatJSON(event.Payload)
			var eventData Transaction
			err := json.Unmarshal(event.Payload, &eventData)
			if err != nil {
				log.Println(err.Error())
			}
			switch event.EventName {
			case "SellVehicle":
				db.InsertTemporalData(eventData.toTemporalData())
				db.InsertSpatialData(eventData.toSpatialData())
				db.InsertParticipantData(eventData.toParticipantData())
				db.InsertTransactionData(eventData.toTransactionData())
				break
			case "BuyVehicle":
				db.UpdateTemporalData(eventData.toTemporalData())
				db.UpdateSpatialData(eventData.toSpatialData())
				db.UpdateParticipantData(eventData.toParticipantData())
				db.UpdateTransactionData(eventData.toTransactionData())
				break
			case "CompromiseTransaction":
				db.UpdateTemporalData(eventData.toTemporalData())
				db.UpdateSpatialData(eventData.toSpatialData())
				db.UpdateParticipantData(eventData.toParticipantData())
				db.UpdateTransactionData(eventData.toTransactionData())
				break
			default:
				fmt.Printf(event.EventName)
			}
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, eventStr)

			endTime := time.Now().UnixNano()
			endTimeString := fmt.Sprintf("%d", endTime)
			if _, err := file.WriteString(startTimeString + " " + endTimeString + "\n"); err != nil {
				log.Println(err)
			}
		}

	}()
}
```
---

### 실시 예
|정보||
|---|---|
|<p style="font-size:10pt" align="center">판매 명세서</p><img width=150/>|<img src=https://github.com/jhl8109/FabricAPI/assets/78259314/1405fbf3-1e77-4122-9096-17b53cf0b7ed width=500 />|
|<p align="center">이벤트 수신 로그</p>|<img src=https://github.com/jhl8109/FabricAPI/assets/78259314/16b7f9e9-e892-4f66-b322-c22dd218235d width=500/>|
|<p align="center">이벤트 처리 로그</p>|<img src=https://github.com/jhl8109/FabricAPI/assets/78259314/a5cde3f6-9694-4db4-9f71-0b3a99d97a50 width=700/>|
|<p align="center">분류별 검색 필터 대시보드</p>|<img src=https://github.com/jhl8109/FabricAPI/assets/78259314/d73d14c1-bff9-4c8e-bf68-8d04f481f8d6 width=700, height=500/>|
<br>

## 성능 평가
> 본 프로젝트에서는 하이퍼레저 패브릭 네트워크와 연결하기 위해 Fabric Gateway SDK를 활용하였으며,<br>
명령어 기반 실행과 SDK 기반 트랜잭션 성능을 평가하였다.

#### 테스트 방법
- 비교 방법 : 실행 속도
- 초기 블록 높이 : 5
- 실행 횟수 : 100회
- 데이터는 쉘스크립트에서 파이프라인을 통해 기록하고 이를 쉘스크립트를 통해 min, 25%, avg, 75%, max 로 종합, 정리함.
- 기존
  - 기존 체인코드는 invoke 시 이벤트를 발생시키지 않도록 구현
  - query 시 단일 데이터에 대해 key기반 조회 수행
- 제안
  - 제안 체인코드는 invoke 시 이벤트를 발생시키도록 구현
  - query 시 단일 데이터에 대해 key기반 조회 수행
  - process는 이벤트로 발생한 데이터를 처리하는 시간을 측정
  

#### 특이사항 
- CLI 첫 값의 latency가 매우 큼
  - 아마 connection 문제일 듯, Gateway는 connection pool이 존재함.
  - 99% line의 값과 maximum값의 차이가 매우 큼, 이는 한 값이 매우 튐을 추측할 수 있음.
  - => 결과적으로 첫 번째 튀는 값을 제외한 나머지들을 통해서 min,25%, avg, 75%, max 값을 평가함.

### 표

|  | 제안 invoke | 제안 process | 제안 query | 기존 invoke | 기존 query |
| --- | --- | --- | --- | --- | --- |
| Minimum | 114.671 | 2.634 | 0.0013 | 107.363 | 75.317 |
| 25% | 140.171 | 4.624 | 0.0015 | 137.604 | 104.527 |
| Average | 158.730 | 7.562 | 0.0020 | 154.359 | 118.426 |
| 75% | 169.571 | 8.575 | 0.0022 | 168.678 | 128.319 |
| Maximum | 214.056 | 60.633 | 0.0074 | 232.173 | 234.933 |

단위 : ms

## 결과
- invoke시 약 12ms만큼 추가 시간이 소요, 그러나, query 시 약 **110ms만큼 성능을 향상** 시킬 수 있음
- 이는 블록체인의 체인 형태로 linked-list와 연결을 타고 올라가는 구조가 아닌 외부 오프체인 데이터베이스의 조회 메커니즘을 활용할 수 있기 때문임(MySQL, 인덱스 기반 조회)
- 따라서, 트랜잭션보다 조회가 자주 일어나는 일반적인 시스템에 더욱 큰 성능 이점을 가질 수 있으며 블록체인 네트워크에 추가적인 부하를 발생시키지 않음.
- 또한, 블록체인 네트워크의 실행환경에서 발생하는 이벤트를 기반으로 동작하므로 온체인 데이터 - 오프체인 데이터의 동기화에 이점이 있음.
- 다만, 데이터에 대한 무결성, 신뢰성, 불변성에 대한 부분이 떨어질 수 있으나 오프체인의 모든 데이터가 온체인에 동기화 되어 있기 때문에 교차 검증을 통해 무결성, 신뢰성을 검증할 수 있음.


