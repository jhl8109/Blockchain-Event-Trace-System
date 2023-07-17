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
- SDK
  - 플랫폼과 연결하는 시간은 포함되지 않음.
  - REST API 요청 ->응답 사이의 시간을 측정하였음.
  - Jmeter를 활용
- CLI
  - 시간 측정은 서버에서 CLI를 통해 쉘스크립트에 작성된 명령어 set을 실행하는 것으로 소요 시간 측정
  - 반복문-쉘스크립트 -> 동작-쉘스크립트 형태로 반복 수행하였음.
  - 데이터는 파이프라인을 통해 기록하고 이를 쉘스크립트를 통해 min,avg,max 로 종합, 정리함.

#### 특이사항
- CLI 첫 값의 latency가 매우 큼
  - 아마 connection 문제일 듯, Gateway는 connection pool이 존재함.
  - 99% line의 값과 maximum값의 차이가 매우 큼, 이는 한 값이 매우 튐을 추측할 수 있음.
  - => 결과적으로 첫 번째 튀는 값을 제외한 나머지들을 통해서 min,avg,max 값을 평가함.

### Charts
| | |
|---|---|
|![img1](https://user-images.githubusercontent.com/78259314/230723374-26c2b3e4-9c85-409f-94bc-78ec8fea9010.png)|![img2](https://user-images.githubusercontent.com/78259314/230723436-cb8fa374-dc61-417e-9c9c-4d26c184e6b9.png)|
|<p align="center">전체 비교</p>|<p align="center">최대</p>|
|![img3](https://user-images.githubusercontent.com/78259314/230723533-4070e3ba-3ed0-4768-8938-afb6b3928e4c.png)|![img4](https://user-images.githubusercontent.com/78259314/230723537-37b80b56-503f-483a-82cb-57853cca28da.png)|
|<p align="center">평균</p>|<p align="center">최소</p>|
<br>

## 결과
- SDK가 CLI보다 성능이 뛰어남.
- 이유 
  - Fabric Gateway SDK의 경우 connection pooling 방식을 통해 패브릭 네트워크와 연결되어있음.
  - 그러나, CLI의 경우 매번 독립적으로 시행되기 때문에 매번 패브릭 네트워크와 연결하는 것으로 추측


